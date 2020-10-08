<?php

declare(strict_types=1);

namespace Tests\Unit\Services\Gateway\Contracts;

use App\Obada\Obit;
use App\Obada\ObitId;
use Tests\TestCase;
use Throwable;
use Illuminate\Validation\ValidationException;

class ObitIdTest extends TestCase {
    /**
     * @test
     */
    public function it_creates_valid_obada_did() {
        $serialHash = hash('sha256', "serial_number");

        $obitId = new ObitId($serialHash, 'manufacturer', 'part number');

        $this->assertEquals(
            'did:obada:8a53fa0b040e11d71d65554e4a887cad14335f7120345c594af794b2e2808184',
            $obitId->toDid()
        );
    }

    /**
     * @test
     */
    public function it_creates_valid_obada_usn() {
        $serialHash = hash('sha256', "serial_number");

        $obitId = new ObitId($serialHash, 'manufacturer', 'part number');

        $this->assertEquals(
            '28NwRR9G',
            $obitId->toUsn()
        );
    }

    /**
     * @test
     */
    public function it_creates_valid_obada_obit() {
        $serialHash = hash('sha256', "serial_number");

        $obitId = new ObitId($serialHash, 'manufacturer', 'part number');

        $this->assertEquals(
            '8a53fa0b040e11d71d65554e4a887cad14335f7120345c594af794b2e2808184',
            $obitId->toHash()
        );
    }
}
