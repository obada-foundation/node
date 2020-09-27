<?php

declare(strict_types=1);

namespace App\Http\Handlers\Obit;

use App\Http\Handlers\Handler;
use App\Services\Gateway\ServiceContract;

class History extends Handler {

    protected ServiceContract $service;

    public function __construct(ServiceContract $service) {
        $this->service = $service;
    }

    /**
     * @param $obitId
     * @return mixed
     */
    public function __invoke($obitId) {
        return $this->service->history($obitId);
    }
}
