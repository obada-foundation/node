<?php

declare(strict_types=1);

namespace App\Services\Blockchain;

use App\Services\Blockchain\Contracts\ServiceContract;
use App\Services\Blockchain\Ion;
use App\Services\Blockchain\Service;
use Aws\QLDB\QLDBClient;
use Aws\QLDBSession\QLDBSessionClient;
use Illuminate\Support\ServiceProvider;

class BlockchainServiceProvider extends ServiceProvider {

    public function register() {
        $this->app->singleton(QLDBClient::class, function() {
            $qldb = new QLDBClient(config('qldb.connection'));

            return $qldb;
        });

        $this->app->singleton(QLDBSessionClient::class, function() {
            $session = new QLDBSessionClient(config('qldb.connection'));

            return $session;
        });

        $this->app->bind(ServiceContract::class, Service::class);
    }
}
