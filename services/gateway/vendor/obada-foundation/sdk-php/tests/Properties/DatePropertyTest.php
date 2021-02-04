<?php

declare(strict_types=1);

namespace Tests\Properties;

use Obada\Exceptions\PropertyValidationException;
use Obada\Properties\DateProperty;
use PHPUnit\Framework\TestCase;

class DatePropertyTest extends TestCase {
    public function testItThrowsValidationExceptionIsValidMethodFails() {
        $this->expectException(PropertyValidationException::class);
        $this->expectExceptionMessage('');

        $this->getMockForAbstractClass(DateProperty::class, [new \DateTime('now')]);
    }
}