<?php

namespace App\Providers;

use Laravel\Lumen\Providers\EventServiceProvider as ServiceProvider;
use App\Services\Gateway\Events\RecordCreated;
use App\Services\Blockchain\Events\RecordCreated as BlockchainRecordCreated;
use App\Listeners\GatewayRecordCreateListener;
use App\Listeners\BlockchainRecordCreateListener;

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
        BlockchainRecordCreated::class => [
            BlockchainRecordCreateListener::class
        ]
    ];
}
