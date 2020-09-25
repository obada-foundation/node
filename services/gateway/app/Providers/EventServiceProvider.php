<?php

namespace App\Providers;

use Laravel\Lumen\Providers\EventServiceProvider as ServiceProvider;
use App\Services\Gateway\Events\RecordCreated;
use App\Listeners\GatewayRecordCreateListener;

class EventServiceProvider extends ServiceProvider
{
    /**
     * The event listener mappings for the application.
     *
     * @var array
     */
    protected $listen = [
        RecordCreated::class => [
            GatewayRecordCreateListener::class,
        ],
    ];
}
