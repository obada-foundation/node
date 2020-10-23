<?php

declare(strict_types=1);

namespace App\Services\Blockchain\QLDB;

use GuzzleHttp\Client;
use Throwable;
use Illuminate\Support\Facades\Log;

class Driver  {

    protected Client $client;

    public function __construct() {
        $this->client = new Client(['base_uri' => 'http://qldb:3000/v1/']);
    }

    /**
     * @param array $obit
     * @return array
     * @throws Throwable
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function create(array $obit) {
        try {
            $response = $this->client->post(
                'obits',
                ['json' => [
                    'obit_did'           => $obit['obit_did'],
                    'usn'                => $obit['usn'],
                    'obit_did_versions'  => '',
                    'owner_did'          => $obit['owner_did'],
                    'obd_did'            => $obit['obd_did'],
                    'serial_number_hash' => $obit['serial_number_hash'],
                    'part_number'        => $obit['part_number'],
                    'manufacturer'       => $obit['manufacturer'],
                    'root_hash'          => $obit['root_hash'],
                    'obit_status'        => $obit['obit_status'],
                    'modified_at'        => $obit['modified_at'],
                    'metadata'           => $obit['metadata'],
                    'doc_links'          => $obit['doc_links'],
                    'structured_data'    => $obit['structured_data']
                ]]
            );
        } catch (Throwable $t) {
            Log::error($t->getMessage(), [$t]);
            throw $t;
        }
    }

    /**
     * @param string $obitId
     * @param array $obit
     * @throws Throwable
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function update(string $obitId, array $obit) {
        try {
            $response = $this->client->put(
                'obits/' . $obitId,
                ['json' => [
                    'obit_did_versions'  => '',
                    'owner_did'          => $obit['owner_did'],
                    'obd_did'            => $obit['obd_did'],
                    'serial_number_hash' => $obit['serial_number_hash'],
                    'part_number'        => $obit['part_number'],
                    'manufacturer'       => $obit['manufacturer'],
                    'root_hash'          => $obit['root_hash'],
                    'obit_status'        => $obit['obit_status'],
                    'modified_at'        => $obit['modified_at'],
                    'metadata'           => $obit['metadata'],
                    'doc_links'          => $obit['doc_links'],
                    'structured_data'    => $obit['structured_data']
                ]]
            );
        } catch (Throwable $t) {
            Log::error($t->getMessage(), [$t]);
            throw $t;
        }
    }

    /**
     * @param string $obitId
     * @throws Throwable
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function delete(string $obitId) {
        try {
            $this->client->delete('obits/' . $obitId);
        } catch (Throwable $t) {
            Log::error($t->getMessage(), [$t]);
            throw $t;
        }
    }
}
