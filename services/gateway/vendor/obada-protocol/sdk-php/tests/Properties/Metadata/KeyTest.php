<?php

declare(strict_types=1);

namespace Tests\Properties\Metadata;

use Obada\Exceptions\PropertyValidationException;
use Obada\Properties\Metadata\Key;
use PHPUnit\Framework\TestCase;

class KeyTest extends TestCase {

	public function testItCreatesMetadataKeyProperty(): void {
		$this->assertEquals('color', (string) new Key('color'));
	}

	public function testItThrowsValidationExceptionWhenMetadataKeyIsEmpty() {
		$this->expectException(PropertyValidationException::class);
		$this->expectExceptionMessage('Key must be valid string with length more than 0.');

		new Key('');
	}

	public function testItGeneratesCorrectHash() {
		$key = new Key('color');

		$this->assertEquals(hash('sha256', 'color'), (string) $key->toHash());
	}
}