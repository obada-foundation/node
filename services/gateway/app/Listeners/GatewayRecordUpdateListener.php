<?php

namespace App\Listeners;

use App\Services\Gateway\Events\RecordUpdated;
use Illuminate\Contracts\Queue\ShouldQueue;
use App\Services\Blockchain\ServiceContract;

class GatewayRecordUpdateListener implements ShouldQueue {
    protected ServiceContract $service;

    /**
     * GatewayRecordCreateListener constructor.
     * @param ServiceContract $service
     */
    public function __construct(ServiceContract $service) {
        $this->service = $service;
    }

    /**
     * @param RecordUpdated $event
     */
    public function handle(RecordUpdated $event) {
        $this->service->update($event->obit->obit_did, $event->obit->toArray());
    }
}
