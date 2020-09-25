<?php

declare(strict_types=1);

namespace App\Http\Handlers\Obit;

use App\Http\Handlers\Handler;
use App\Services\Gateway\Contracts\ServiceContract;
use App\Http\Requests\Obit\CreateRequest;
use App\Services\Gateway\ObitDto;

class Create extends Handler {

    /**
     * @var ServiceContract
     */
    protected $service;

    /**
     * Create constructor.
     * @param ServiceContract $service
     */
    public function __construct(ServiceContract $service) {
        $this->service = $service;
    }

    /**
     * @param CreateRequest $request
     * @return \Illuminate\Http\Response
     */
    public function __invoke(CreateRequest $request) {
        $dto = ObitDto::fromRequest($request);

        $this->service->create($dto);

        return $this->successRequestWithNoData();
    }
}
