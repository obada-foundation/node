<?php

namespace App\Listeners;

use App\Services\Gateway\Events\RecordCreated;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Queue\InteractsWithQueue;
use App\Services\Blockchain\Contracts\ServiceContract;

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
     * Handle the event.
     *
     * @param  \App\Events\RecordCreated  $event
     * @return void
     */
    public function handle(RecordCreated $event)
    {
        $this->service->create($event->obit->toArray());
    }
}
