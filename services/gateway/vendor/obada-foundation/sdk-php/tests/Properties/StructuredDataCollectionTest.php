<?php

declare(strict_types=1);

namespace Tests\Properties;

use PHPUnit\Framework\TestCase;
use Obada\Properties\StructuredData\Record;
use Obada\Properties\StructuredData\Key;
use Obada\Properties\StructuredData\Value;
use Obada\Properties\StructuredDataCollection;

class StructuredDataCollectionTest extends TestCase {

    public function testItCreatesStructuredDataCollection() {
        $collection = new StructuredDataCollection;

        $this->assertCount(0, $collection);
        $this->assertEquals([], $collection->toArray());

        $record = new Record(new Key('foo'), new Value('bar'));

        $collection->add($record);

        $this->assertCount(1, $collection);

        $this->assertCount(1, new StructuredDataCollection($record));

        $this->assertCount(2, new StructuredDataCollection($record, $record));
    }
}