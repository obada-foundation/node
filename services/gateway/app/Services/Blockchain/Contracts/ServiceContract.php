<?php

declare(strict_types=1);

namespace App\Services\Blockchain\Contracts;

interface ServiceContract {
    public function create(array $obit);

    public function show(string $obitId);

    public function update(string $obitId);

    public function delete(string $obitId);

    public function history(string $obitId);
}
