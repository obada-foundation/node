<?php

declare(strict_types=1);

namespace App\Http\Handlers\Obit;

use App\Http\Handlers\Handler;
use App\Services\Gateway\Contracts\ServiceContract;

class Search extends Handler {

    /**
     * @var ServiceContract
     */
    protected $service;

    /**
     * Search constructor.
     * @param ServiceContract $service
     */
    public function __construct(ServiceContract $service) {
        $this->service = $service;
    }

    /**
     * @return mixed
     */
    public function __invoke() {
        $obits = $this->service->search();

        return $obits;
    }
}
