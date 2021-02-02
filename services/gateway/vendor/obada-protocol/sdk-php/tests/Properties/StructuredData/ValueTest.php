<?php

declare(strict_types=1);

namespace Tests\Properties\StructuredData;

use Obada\Exceptions\PropertyValidationException;
use Obada\Properties\StructuredData\Value;
use PHPUnit\Framework\TestCase;

class ValueTest extends TestCase {

    public function testItCreatesStructuredDataKey() {
        $this->assertEquals('bar', new Value('bar'));
    }

    public function testItThrowsValidationExceptionWhenStructuredDateKeyIsNotValid() {
        $this->expectException(PropertyValidationException::class);
        $this->expectExceptionMessage('Value must be valid string with length more than 0.');

        new Value('');
    }
}