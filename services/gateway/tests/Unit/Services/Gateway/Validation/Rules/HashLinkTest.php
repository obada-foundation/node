<?php

declare(strict_types=1);

namespace Tests\Unit\Services\Gateway\Validation\Rules;

use App\Services\Gateway\Validation\Rules\HashLink;
use Illuminate\Support\Facades\Validator;
use Illuminate\Validation\ValidationException;
use Tests\TestCase;
use Throwable;

class HashLinkTest extends TestCase {
    /**
     * @test
     */
    public function it_fails_where_no_hl_query_argument_was_provided() {
        try {
            Validator::make(
                ['hl' => 'https://google.com'],
                ['hl' => new HashLink]
            )->validate();

            $this->assertTrue(false);
        } catch (Throwable $t) {
            $this->assertInstanceOf(ValidationException::class, $t);
            $this->assertEquals($t->errors(), ['hl' => ['Not a valid hashlink.']]);
        }
    }

    /**
     * @test
     */
    public function it_fails_where_no_hl_query_argument_exists_but_null() {
        try {
            Validator::make(
                ['hl' => 'https://google.com?hl'],
                ['hl' => new HashLink]
            )->validate();

            $this->assertTrue(false);
        } catch (Throwable $t) {
            $this->assertInstanceOf(ValidationException::class, $t);
            $this->assertEquals($t->errors(), ['hl' => ['Not a valid hashlink.']]);
        }
    }

    /**
     * @test
     */
    public function it_fails_where_no_hl_query_argument_exists_but_url_doesnt_provide_data() {
        try {
            Validator::make(
                ['hl' => 'https://go2o2gle.com?hl=123'],
                ['hl' => new HashLink]
            )->validate();

            $this->assertTrue(false);
        } catch (Throwable $t) {
            $this->assertInstanceOf(ValidationException::class, $t);
            $this->assertEquals($t->errors(), ['hl' => ['Not a valid hashlink.']]);
        }
    }

    /**
     * @test
     */
    public function it_fails_where_document_hash_and_hashlink_doesnt_match() {
        try {
            Validator::make(
                ['hl' => 'https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf?hl=123'],
                ['hl' => new HashLink]
            )->validate();

            $this->assertTrue(false);
        } catch (Throwable $t) {
            $this->assertInstanceOf(ValidationException::class, $t);
            $this->assertEquals($t->errors(), ['hl' => ['Not a valid hashlink.']]);
        }
    }

    /**
     * @test
     */
    public function it_passes_where_document_hash_and_hashlink_are_match() {
        $hash = hash('sha256', file_get_contents(
            'https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf')
        );

        Validator::make(
            ['hl' => 'https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf?hl=' . $hash],
            ['hl' => new HashLink]
        )->validate();

        $this->assertTrue(true);
    }
}
