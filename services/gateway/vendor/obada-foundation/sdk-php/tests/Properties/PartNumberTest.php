<?php

declare(strict_types=1);

namespace Tests\Properties;

use Obada\Exceptions\PropertyValidationException;
use Obada\Properties\PartNumber;
use PHPUnit\Framework\TestCase;

class PartNumberTest extends TestCase {

	public function testItCreatesPartNumberProperty(): void {
		$this->assertEquals("part number", new PartNumber('part number'));
	}

	public function testItThrowsValidationExceptionWhenPartNumberIsEmpty() {
	    $this->expectException(PropertyValidationException::class);
	    $this->expectExceptionMessage('PartNumber is required and cannot be empty');

	    new PartNumber('');
	}
}