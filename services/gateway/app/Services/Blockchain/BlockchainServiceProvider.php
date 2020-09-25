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
        $this->app->bind(ServiceContract::class, Service::class);
    }
}
