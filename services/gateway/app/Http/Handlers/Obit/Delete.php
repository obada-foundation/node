<?php

declare(strict_types=1);

namespace App\Http\Handlers\Obit;

use App\Http\Handlers\Handler;
use App\Services\Gateway\ServiceContract;

class Delete extends Handler {

    protected ServiceContract $service;

    public function __construct(ServiceContract $service) {
        $this->service = $service;
    }

    /**
     * @return \Illuminate\Http\Response
     */
    public function __invoke() {
        $did = request()->route()[2]['obitDID'];

        $this->service->delete($did);

        return $this->successRequestWithNoData();
    }
}
