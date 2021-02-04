<?php

declare(strict_types=1);

namespace Tests\Properties;

use PHPUnit\Framework\TestCase;
use Obada\Properties\Document\Document;
use Obada\Properties\Document\Name;
use Obada\Properties\Document\HashLink;
use Obada\Properties\DocumentsCollection;

class DocumentsCollectionTest extends TestCase {

    public function testItCreatesDocumentsCollection() {
        $collection = new DocumentsCollection;

        $this->assertCount(0, $collection);
        $this->assertEquals([], $collection->toArray());

        $document = new Document(new Name('swipe report'), new HashLink('http://some-link.com'));

        $collection->add($document);

        $this->assertCount(1, $collection);

        $this->assertCount(1, new DocumentsCollection($document));

        $this->assertEquals([['name' => 'swipe report', 'hash_link' => 'http://some-link.com']], $collection->toArray());

        $collection = new DocumentsCollection($document, $document);

        $this->assertCount(2, $collection);

        $this->assertEquals(
            [
                ['name' => 'swipe report', 'hash_link' => 'http://some-link.com'],
                ['name' => 'swipe report', 'hash_link' => 'http://some-link.com']
            ],
            $collection->toArray()
        );

    }
}