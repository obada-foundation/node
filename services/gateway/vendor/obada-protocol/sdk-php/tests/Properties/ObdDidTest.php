<?php

declare(strict_types=1);

namespace Tests\Properties;

use Obada\Properties\ObdDid;
use PHPUnit\Framework\TestCase;

class ObdDidTest extends TestCase {
    public function testItCreatesObdDidProperty(): void {
        $this->assertEquals('did:obada:obd:1234', new ObdDid('did:obada:obd:1234'));
    }
}