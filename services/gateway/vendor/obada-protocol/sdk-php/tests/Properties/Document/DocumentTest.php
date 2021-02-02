<?php

declare(strict_types=1);

namespace Tests\Properties\Document;

use Obada\Properties\Document\Name;
use Obada\Properties\Document\Document;
use Obada\Properties\Document\HashLink;
use PHPUnit\Framework\TestCase;
use Obada\Hash;

class DocumentTest extends TestCase {

    public function testItCreatesDocumentProperty(): void {
        $name = new Name('swipe report');
        $link = new HashLink('some link');

        $record = new Document($name, $link);

        $this->assertEquals($name, $record->getName());
        $this->assertEquals($link, $record->getHashLink());

        $expectedHash = new Hash(dechex($name->toHash()->toDecimal() + $link->toHash()->toDecimal()));

        $this->assertEquals($expectedHash, $record->toHash());
    }
}