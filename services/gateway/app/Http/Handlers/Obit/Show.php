<?php

declare(strict_types=1);

namespace App\Http\Handlers\Obit;

use App\Http\Handlers\Handler;
use App\Services\Gateway\ServiceContract;

class Show extends Handler {

    protected ServiceContract $service;

    /**
     * Show constructor.
     * @param ServiceContract $service
     */
    public function __construct(ServiceContract $service) {
        $this->service = $service;
    }

    /**
     * @param string $obitDID
     * @return \Illuminate\Http\JsonResponse
     */
    public function __invoke() {
        $did  = request()->route()[2]['obitDID'];
        $obit = $this->service->show($did);

        if (! $obit) {
            return $this->responseNotFound("Obit with did \"{$did}\" doesn't exists");
        }

        return $this->respond([
            'obit_did'           => (string) $obit->obit_did,
            'usn'                => (string) $obit->usn,
            'obit_did_versions'  => (array) $obit->obit_did_versions,
            'owner_did'          => (string) $obit->owner_did,
            'obd_did'            => (string) $obit->obd_did,
            'obit_status'        => (string) $obit->obit_status,
            'manufacturer'       => (string) $obit->manufacturer,
            'part_number'        => (string) $obit->part_number,
            'serial_number_hash' => (string) $obit->serial_number_hash,
            'metadata'           => (array) $obit->metadata,
            'structured_data'    => (array) $obit->structured_data,
            'doc_links'          => (array) $obit->doc_links,
            'modified_at'        => (string) $obit->modified_at,
            'root_hash'          => (string) $obit->root_hash,
            'qldb_root_hash'     => (string) $obit->qldb_root_hash,
            'is_synchronized'    => (boolean) $obit->is_synchronized
        ]);
    }
}
