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
        $query = Obit::query();

        if (isset($args['query']) && $q = $args['query']) {
            $query->where('obit_did', $q)
                ->orWhere('usn', $q)
                ->orWhere('manufacturer', $q)
                ->orWhere('part_number', $q)
                ->orWhere('serial_number_hash', $q);;
        }

        if (isset($args['serial_number_hash']) && $serialNumberHash = $args['serial_number_hash']) {
            $query->where('serial_number_hash', $serialNumberHash);;
        }

        if (isset($args['obit_status']) && $obitStatus = $args['obit_status']) {
            $query->where('obit_status', $obitStatus);;
        }

        if (isset($args['manufacturer']) && $manufacturer = $args['manufacturer']) {
            $query->where('manufacturer', $manufacturer);;
        }

        if (isset($args['part_number']) && $partNumber = $args['part_number']) {
            $query->where('part_number', $partNumber);;
        }

        if (isset($args['usn']) && $usn = $args['usn']) {
            $query->where('usn', $usn);;
        }

        if (isset($args['owner_did']) && $ownerDID = $args['owner_did']) {
            $query->where('owner_did', $ownerDID);;
        }

        return $query->paginate(50);
    }

    /**
     * @param string $obitDID
     * @return Obit
     */
    public function find(string $obitDID): ?Obit {
        return Obit::where('obit_did', $obitDID)->first();
    }
}
