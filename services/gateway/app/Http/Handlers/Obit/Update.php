<?php

declare(strict_types=1);

namespace App\Http\Handlers\Obit;

use App\Http\Handlers\Handler;
use App\Obada\Mappers\Input\ObitInputMapper;
use App\Services\Gateway\UpdateObitDto;
use App\Services\Gateway\ServiceContract;
use Illuminate\Support\Facades\Log;
use Obada\Exceptions\PropertyValidationException;

class Update extends Handler {

    protected ServiceContract $service;

    public function __construct(ServiceContract $service) {
        $this->service = $service;
    }

    /**
     * @return \Illuminate\Http\Response
     */
    public function __invoke() {
        Log::debug('request', [request()->all()]);

        $did = request()->route()[2]['obitDID'];

        try {
            $obit = app()
                ->make(ObitInputMapper::class)
                ->map(request()->json()->all());

            $this->service->update($did, $obit);
        } catch (PropertyValidationException $t) {
            return $this->respondValidationErrors([], $t->getMessage());
        }

        return $this->successRequestWithNoData();
    }
}
