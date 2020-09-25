<?php

declare(strict_types=1);

namespace Tests\Unit\Services\Gateway\Contracts;

use App\Services\Gateway\Contracts\ServiceContract;
use Tests\TestCase;

class ServiceContractTest extends TestCase {

    protected $service;

    public function setUp(): void {
        parent::setUp(); // TODO: Change the autogenerated stub

        $this->service = app()->make(ServiceContract::class);
    }

    /**
     * @test
     */
    public function it_creates_a_view_record() {
        //$this->service->create();
    }
}