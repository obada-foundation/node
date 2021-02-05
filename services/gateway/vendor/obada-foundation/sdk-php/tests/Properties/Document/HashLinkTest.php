<?php

declare(strict_types=1);

namespace Tests\Properties\Document;

use Obada\Exceptions\PropertyValidationException;
use Obada\Properties\Document\HashLink;
use PHPUnit\Framework\TestCase;

class HashLinkTest extends TestCase {

    public function testItCreatesDocumentNameProperty(): void {
        $this->assertEquals('swipe report', (string) new HashLink('swipe report'));
    }

    /**
     * @dataProvider getInvalidDocumentNames
     */
    public function testItThrowsValidationExceptionWhenDocumentNameIsEmpty($documentName) {
        $this->expectException(PropertyValidationException::class);
        $this->expectExceptionMessage('Hash link must be valid string with length more than 0.');

        new HashLink($documentName);
    }

    public function getInvalidDocumentNames(): \Generator {
        yield [''];
    }

    public function testItGeneratesCorrectHash() {
        $name = new HashLink('swipe report');

        $this->assertEquals(hash('sha256', 'swipe report'), $name->toHash());
    }
}