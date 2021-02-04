<?php

declare(strict_types=1);

namespace Tests\Properties\Document;

use Obada\Exceptions\PropertyValidationException;
use Obada\Properties\Document\Name;
use PHPUnit\Framework\TestCase;

class NameTest extends TestCase {

    public function testItCreatesDocumentNameProperty(): void {
        $this->assertEquals('swipe report', (string) new Name('swipe report'));
    }

    /**
     * @dataProvider getInvalidDocumentNames
     */
    public function testItThrowsValidationExceptionWhenDocumentNameIsEmpty($documentName) {
        $this->expectException(PropertyValidationException::class);
        $this->expectExceptionMessage('Name must be not empty string with length more than 0.');

        new Name($documentName);
    }

    public function getInvalidDocumentNames(): \Generator {
        yield [''];
        yield [' '];
    }

    public function testItGeneratesCorrectHash() {
        $name = new Name('swipe report');

        $this->assertEquals(hash('sha256', 'swipe report'), $name->toHash());
    }
}