<?php

declare(strict_types=1);

namespace Tests\Unit\Services\Gateway;

use Illuminate\Validation\ValidationException;
use Tests\TestCase;
use App\Services\Gateway\ObitDto;
use Throwable;

class ObitDtoTest extends TestCase {

    protected string $hash = '8a53fa0b040e11d71d65554e4a887cad14335f7120345c594af794b2e2808184';

    protected array $dto = [
        'modifiedAt'       => '20-03-2020 ',
        'serialNumberHash' => '8a53fa0b040e11d71d65554e4a887cad14335f7120345c594af794b2e2808184',
        'manufacturer'     => 'manufacturer',
        'usn'              => '28NwRR9G',
        'partNumber'       => 'part number',
        'obitDID'          => '8a53fa0b040e11d71d65554e4a887cad14335f7120345c594af794b2e2808184'
    ];

    /**
     * @test
     */
    public function it_doesnt_pass_validation_when_no_required_fields_were_passed() {
        try {
            new ObitDto();

            $this->assertTrue(false);
        } catch (Throwable $t) {
            $this->assertInstanceOf(ValidationException::class, $t);
            $this->assertEquals($t->errors(), ['serial_number_hash' => ['The serial number hash field is required.']]);
        }

        try {
            new ObitDto(['serialNumberHash' => 'serial']);

            $this->assertTrue(false);
        } catch (Throwable $t) {
            $this->assertInstanceOf(ValidationException::class, $t);
            $this->assertEquals(
                $t->errors(),
                ['serial_number_hash' => ['Received the invalid serial number hash. Must be the valid SHA256 hash.']]
            );
        }

        try {
            new ObitDto(['serialNumberHash' => $this->hash, 'obitDID' => '']);

            $this->assertTrue(false);
        } catch (Throwable $t) {
            $this->assertInstanceOf(ValidationException::class, $t);
            $this->assertEquals(
                $t->errors(),
                [
                    'obit_did'     => ['The obit did field is required.'],
                    'manufacturer' => ['The manufacturer field is required.'],
                    'usn'          => ['The usn field is required.'],
                    'modified_at'  => ['The modified at field is required.'],
                    'part_number'  => ['The part number field is required.'],
                ]
            );
        }
    }

    /**
     * @test
     */
    public function it_doest_pass_validation_when_obit_id_has_integrity_check_problems() {
        try {
            $dto = $this->dto;
            $dto['manufacturer'] = 'man1';

            new ObitDto($this->dto);

            $this->assertTrue(false);
        } catch (Throwable $t) {
            $this->assertInstanceOf(ValidationException::class, $t);
            $this->assertEquals(
                $t->errors(),
                ['usn' => ['The selected usn is invalid.'], 'obit_did' => ['Integrity of obit id is broken.']]
            );
        }
    }
}
