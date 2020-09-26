<?php

namespace App\Listeners;

use App\Services\Gateway\Events\RecordCreated;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Queue\InteractsWithQueue;
use App\Services\Blockchain\ServiceContract;

class GatewayRecordCreateListener
{
    protected ServiceContract $service;

    /**
     * GatewayRecordCreateListener constructor.
     * @param ServiceContract $service
     */
    public function __construct(ServiceContract $service)
    {
        $this->service = $service;
    }

    /**
     * @param RecordCreated $event
     */
    public function handle(RecordCreated $event)
    {
        $this->service->create($event->obit->toArray());
    }
}
