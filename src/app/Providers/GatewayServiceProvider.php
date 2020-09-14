<?php

declare(strict_types=1);

namespace App\Providers;

use App\Services\Gateway\Contracts\ServiceContract;
use App\Services\Gateway\Service;
use Illuminate\Support\ServiceProvider;

class GatewayServiceProvider extends ServiceProvider {

    public function register() {
        $this->app->bind(ServiceContract::class, Service::class);
    }
}
