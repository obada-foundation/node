<?php

declare(strict_types=1);

namespace Tests\Properties\Metadata;

use Obada\Hash;
use Obada\Properties\Metadata\Key;
use Obada\Properties\Metadata\Record;
use Obada\Properties\Metadata\Value;
use PHPUnit\Framework\TestCase;

class RecordTest extends TestCase {

	public function testItCreatesMetadataRecordProperty(): void {
		$key   = new Key('color');
		$value = new Value('red');

		$record = new Record($key, $value);

		$this->assertEquals($key, $record->getKey());
		$this->assertEquals($value, $record->getValue());

        $expectedHash = new Hash(dechex($key->toHash()->toDecimal() + $value->toHash()->toDecimal()));

        $this->assertEquals($expectedHash, $record->toHash());
	}
}