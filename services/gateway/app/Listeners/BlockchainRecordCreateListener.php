<?php

namespace App\Listeners;

use App\Services\Blockchain\Events\RecordCreated;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Queue\InteractsWithQueue;
use App\Services\Gateway\ServiceContract;

class BlockchainRecordCreateListener implements ShouldQueue {
    protected ServiceContract $service;

    /**
     * GatewayRecordCreateListener constructor.
     * @param ServiceContract $service
     */
    public function __construct(ServiceContract $service) {
        $this->service = $service;
    }

    /**
     * @param RecordCreated $event
     */
    public function handle(RecordCreated $event) {
        $this->service->commit($event->obit['obit_did']);
    }
}
