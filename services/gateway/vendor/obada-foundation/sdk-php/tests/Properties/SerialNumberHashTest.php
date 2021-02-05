<?php

declare(strict_types=1);

namespace Tests\Properties;

use Obada\Exceptions\PropertyValidationException;
use Obada\Properties\SerialNumberHash;
use PHPUnit\Framework\TestCase;
use \Generator;

class SerialNumberHashTest extends TestCase {
    public function testItCreatesSerialNumberHashProperty(): void {
        $hash = hash('sha256', 'serial number');

        $this->assertEquals($hash, new SerialNumberHash($hash));
    }

    /**
     * @dataProvider getInvalidSerialNumberHashes
     */
    public function testItThrowsValidationExceptionWhenSerialNumberHashIsNotValid(string $invalidSerialNumberHash) {
        $this->expectException(PropertyValidationException::class);
        $this->expectExceptionMessage('Serial number hash must be a valid SHA256 hash');

        new SerialNumberHash($invalidSerialNumberHash);
    }

    function getInvalidSerialNumberHashes(): Generator {
        yield ['1'];
        yield ['string'];
        yield [md5('string')];
    }
}