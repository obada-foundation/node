<?php

declare(strict_types=1);

namespace Tests\Feature\Http\Obit;

use Carbon\Carbon;
use Laravel\Lumen\Testing\DatabaseTransactions;
use Laravel\Lumen\Testing\WithoutEvents;
use Obada\Obit;
use Tests\TestCase;

class CreateTest extends TestCase {

    use DatabaseTransactions, WithoutEvents;

    /**
     * @test
     */
    public function it_validates_that_obit_has_required_fields() {
        $payload = [];

        $this->json("POST", route('obits.create'), $payload);

        $this->seeStatusCode(422);
        $this->seeJson([
            'code'    => 422,
            'errors'  => [],
            'message' => 'Manufacturer is required and cannot be empty'
        ]);

        $payload = [
            'manufacturer' => 'Dell'
        ];

        $this->json("POST", route('obits.create'), $payload);

        $this->seeStatusCode(422);
        $this->seeJson([
            'code'    => 422,
            'errors'  => [],
            'message' => 'Serial number hash must be a valid SHA256 hash'
        ]);

        $payload = [
            'manufacturer' => 'Dell'
        ];

        $this->json("POST", route('obits.create'), $payload);

        $this->seeStatusCode(422);
        $this->seeJson([
            'code'    => 422,
            'errors'  => [],
            'message' => 'Serial number hash must be a valid SHA256 hash'
        ]);
    }

    /**
     * @test
     */
    public function it_returns_correct_response_when_create_basic_obit() {
        $obit = Obit::make([
            'manufacturer'       => 'Sony',
            'serial_number_hash' => hash('sha256', 'SN123456'),
            'part_number'        => 'PN123456',
            'owner_did'          => '123456',
            'modified_at'        => Carbon::now(),
            'obit_status'        => 'FUNCTIONAL',
            'metadata'           => ['color' => 'red'],
            'structured_data'    => ['condition' => 'good']
        ]);

        $metadata = (array) collect($obit->getMetadata()->toArray())
            ->map(fn($record) => [(string) $record->getKey() => (string) $record->getValue()])
            ->collapse()
            ->toArray();

        $structuredData = (array) collect($obit->getStructuredData()->toArray())
            ->map(fn($record) => [(string) $record->getKey() => (string) $record->getValue()])
            ->collapse()
            ->toArray();

        $payload = [
            'manufacturer'       => (string) $obit->getManufacturer(),
            'serial_number_hash' => (string) $obit->getSerialNumberHash(),
            'part_number'        => (string) $obit->getPartNumber(),
            'owner_did'          => (string) $obit->getOwnerDid(),
            'modified_at'        => (string) $obit->getModifiedAt(),
            'root_hash'          => (string) $obit->rootHash(),
            'obit_status'        => (string) $obit->getStatus(),
            'metadata'           => $metadata,
            'structured_data'    => $structuredData
        ];

        $this->json("POST", route('obits.create'), $payload);
        $this->seeStatusCode(204);
    }
}
