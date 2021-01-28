<?php

declare(strict_types=1);

namespace App\Services\Gateway;

use Obada\Obit;
use App\Services\Gateway\Models\Obit as Model;
use Illuminate\Support\Collection;

interface ServiceContract {
    /**
     * @return mixed
     */
    public function search();

    /**
     * @param Obit $obit
     * @return mixed
     */
    public function create(Obit $obit);

    /**
     * @param string $obitId
     * @param Obit $obit
     * @return mixed
     */
    public function update(string $obitId, Obit $obit);

    /**
     * @param string $obitId
     * @return Model|null
     */
    public function show(string $obitId): ?Model;

    /**
     * @param string $obitId
     * @return mixed
     */
    public function delete(string $obitId);

    /**
     * @param string $obitId
     * @return Collection
     */
    public function history(string $obitId): Collection;

    /**
     * @param array $qldbMetadata
     * @return mixed
     */
    public function commit(array $qldbMetadata);
}
