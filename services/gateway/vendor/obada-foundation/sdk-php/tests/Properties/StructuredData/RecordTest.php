<?php

declare(strict_types=1);

namespace Tests\Properties\StructuredData;

use Obada\Hash;
use Obada\Properties\StructuredData\Key;
use Obada\Properties\StructuredData\Record;
use Obada\Properties\StructuredData\Value;
use PHPUnit\Framework\TestCase;

class RecordTest extends TestCase {
    public function testItCreatesStructuredDataRecord() {
        $key    = new Key('foo');
        $value  = new Value('bar');
        $record = new Record($key, $value);

        $this->assertEquals($key, $record->getKey());
        $this->assertEquals($value, $record->getValue());

        $expectedHash = new Hash(dechex($key->toHash()->toDecimal() + $value->toHash()->toDecimal()));

        $this->assertEquals($expectedHash, $record->toHash());
    }
}