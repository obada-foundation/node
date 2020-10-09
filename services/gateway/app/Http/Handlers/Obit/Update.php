<?php

declare(strict_types=1);

namespace App\Http\Handlers\Obit;

use App\Http\Handlers\Handler;
use App\Services\Gateway\UpdateObitDto;
use App\Services\Gateway\ServiceContract;
use App\Http\Requests\Obit\UpdateRequest;
use Illuminate\Support\Facades\Log;

class Update extends Handler {

    protected ServiceContract $service;

    public function __construct(ServiceContract $service) {
        $this->service = $service;
    }

    /**
     * @param UpdateRequest $request
     * @return \Illuminate\Http\Response
     */
    public function __invoke(UpdateRequest $request) {
        Log::debug('request', [$request]);

        $did = request()->route()[2]['obitDID'];

        $dto = UpdateObitDto::fromRequest($request);

        $this->service->update($did, $dto);

        return $this->successRequestWithNoData();
    }
}
