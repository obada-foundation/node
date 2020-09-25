<?php

declare(strict_types=1);

namespace App\Http\Handlers\Obit;

use App\Http\Handlers\Handler;
use App\Services\Gateway\Contracts\ServiceContract;

class History extends Handler {

    protected $service;

    public function __construct(ServiceContract $service) {
        $this->service = $service;
    }

    public function __invoke() {

    }
}
