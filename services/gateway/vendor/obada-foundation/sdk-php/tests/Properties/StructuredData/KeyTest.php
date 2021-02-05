<?php

declare(strict_types=1);

namespace Tests\Properties\StructuredData;

use Obada\Exceptions\PropertyValidationException;
use Obada\Properties\StructuredData\Key;
use PHPUnit\Framework\TestCase;
use \Generator;

class KeyTest extends TestCase {

    public function testItCreatesStructuredDataKey() {
        $this->assertEquals('foo', new Key('foo'));
    }

    /**
     * @dataProvider getInvalidSerialNumberHashes
     */
    public function testItThrowsValidationExceptionWhenStructuredDateKeyIsNotValid($key) {
        $this->expectException(PropertyValidationException::class);
        $this->expectExceptionMessage('Key must be valid string with length more than 0.');

        new Key($key);
    }

    function getInvalidSerialNumberHashes(): Generator {
        yield [''];
        yield [' '];
    }
}