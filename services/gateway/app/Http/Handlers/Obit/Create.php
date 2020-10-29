<?php

declare(strict_types=1);

namespace App\Http\Handlers\Obit;

use App\Http\Handlers\Handler;
use App\Services\Gateway\ServiceContract;
use Illuminate\Support\Facades\Log;
use App\Obada\Mappers\Input\ObitInputMapper;
use Obada\Exceptions\PropertyValidationException;

class Create extends Handler {

    /**
     * @var ServiceContract
     */
    protected ServiceContract $service;

    /**
     * Create constructor.
     * @param ServiceContract $service
     */
    public function __construct(ServiceContract $service) {
        $this->service = $service;
    }

    /**
     * @return \Illuminate\Http\JsonResponse|\Illuminate\Http\Response
     * @throws \Symfony\Component\HttpFoundation\Exception\BadRequestException
     */
    public function __invoke() {
        // Move to middleware
        Log::debug('request', [request()->all()]);

        try {
            $obit = app()
                ->make(ObitInputMapper::class)
                ->map(request()->json()->all());

            $this->service->create($obit);
        } catch (PropertyValidationException $t) {
            return $this->respondValidationErrors([], $t->getMessage());
        }

        return $this->successRequestWithNoData();
    }
}
