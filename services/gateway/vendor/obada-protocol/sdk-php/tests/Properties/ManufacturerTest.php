<?php

declare(strict_types=1);

namespace Tests\Properties;

use Obada\Exceptions\PropertyValidationException;
use Obada\Properties\Manufacturer;
use PHPUnit\Framework\TestCase;

class ManufacturerTest extends TestCase {
    public function testItCreatesManufacturerProperty(): void {
        $this->assertEquals('sony', new Manufacturer('sony'));
    }

    public function testItThrowsValidationExceptionWhenManufacturerIsEmpty() {
        $this->expectException(PropertyValidationException::class);
        $this->expectExceptionMessage('Manufacturer is required and cannot be empty');

        new Manufacturer('');
    }
}