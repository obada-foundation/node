<?php

declare(strict_types=1);

namespace Tests\Properties;

use Obada\Exceptions\PropertyValidationException;
use Obada\Properties\OwnerDid;
use PHPUnit\Framework\TestCase;

class OwnerDidTest extends TestCase {
    public function testItCreatesOwnerDidProperty(): void {
        $this->assertEquals('did:obada:owner:123456', new OwnerDid('did:obada:owner:123456'));
    }

    public function testItThrowsValidationExceptionWhenOwnerDidIsEmpty() {
        $this->expectException(PropertyValidationException::class);
        $this->expectExceptionMessage('OwnerDid is required and cannot be empty');

        new OwnerDid('');
    }
}