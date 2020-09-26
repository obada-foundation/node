<?php

declare(strict_types=1);

namespace App\Services\Blockchain;

use Illuminate\Support\ServiceProvider;

class BlockchainServiceProvider extends ServiceProvider {

    public function register() {
        $this->app->bind(ServiceContract::class, Service::class);
    }
}
