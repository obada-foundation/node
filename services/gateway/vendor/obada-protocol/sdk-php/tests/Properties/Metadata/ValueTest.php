<?php

declare(strict_types=1);

namespace Tests\Properties\Metadata;

use Obada\Exceptions\PropertyValidationException;
use Obada\Properties\Metadata\Value;
use PHPUnit\Framework\TestCase;

class ValueTest extends TestCase {

    public function testItCreatesMetadataValueProperty(): void {
        $this->assertEquals('red', (string) new Value('red'));
    }

    public function testItThrowsValidationExceptionWhenMetadataValueIsEmpty() {
        $this->expectException(PropertyValidationException::class);
        $this->expectExceptionMessage('Value must be valid string with length more than 0.');

        new Value('');
    }

    public function testItGeneratesCorrectHash() {
        $key = new Value('red');

        $this->assertEquals(hash('sha256', 'red'), (string) $key->toHash());
    }
}