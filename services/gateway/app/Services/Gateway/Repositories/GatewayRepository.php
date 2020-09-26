<?php

declare(strict_types=1);

namespace App\Services\Gateway\Repositories;

use App\Services\Gateway\Models\Obit;
use Illuminate\Pagination\LengthAwarePaginator;

class GatewayRepository implements GatewayRepositoryContract {

    /**
     * @param array $args
     * @return Collection
     */
    public function findBy(array $args = []): LengthAwarePaginator {
        return Obit::paginate(50);
    }

    /**
     * @param string $obitDID
     * @return Obit
     */
    public function find(string $obitDID): ?Obit {
        return Obit::where('obit_did', $obitDID)->first();
    }
}
