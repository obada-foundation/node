<?php

declare(strict_types=1);

namespace Tests\Unit\Services\Gateway\Validation\Rules;

use App\Obada\ObitId;
use Illuminate\Support\Facades\Validator;
use Illuminate\Validation\ValidationException;
use Tests\TestCase;
use App\Services\Gateway\Validation\Rules\ObitIntegrity;
use Throwable;

class ObitIntegrityTest extends TestCase {
    /**
     * @test
     */
    public function it_fails_when_obit_integrity_is_there() {
        try {
            $obitId = new ObitId(hash('sha256', 'serial number'), 'manufacturer', 'part number');

            Validator::make(
                ['obit_id' => hash('sha256', '123')],
                ['obit_id' => new ObitIntegrity($obitId)],
            )->validate();

            $this->assertTrue(false);
        } catch (Throwable $t) {
            $this->assertInstanceOf(ValidationException::class, $t);
            $this->assertEquals($t->errors(), ['obit_id' => ['Integrity of obit id is broken.']]);
        }
    }
}
