<?php

declare(strict_types=1);

namespace Tests\Unit\Services\Gateway;

use App\Services\Gateway\Models\Obit;
use App\Services\Gateway\ServiceContract;
use Carbon\Carbon;
use Tests\TestCase;
use Laravel\Lumen\Testing\DatabaseTransactions;

class ServiceContractTest extends TestCase {

    use DatabaseTransactions;

    protected $service;

    public function setUp(): void {
        parent::setUp(); // TODO: Change the autogenerated stub

        $this->service = app()->make(ServiceContract::class);
    }

    /**
     * @test
     */
    public function it_creates_a_view_record() {
        $this->withoutEvents();

        $obit = \Obada\Obit::make([
            'manufacturer'       => 'Sony',
            'serial_number_hash' => hash('sha256', 'SN123456'),
            'part_number'        => 'PN123456',
            'owner_did'          => '123456',
            'modified_at'        => Carbon::now()
        ]);

        $this->service->create($obit);

        $this->assertCount(1, Obit::all());

        $this->seeInDatabase('gateway_view', [
            'obit_did'           => (string) $obit->getObitId()->toDid(),
            'usn'                => (string) $obit->getObitId()->toUsn(),
            'manufacturer'       => (string) $obit->getManufacturer(),
            'part_number'        => (string) $obit->getPartNumber(),
            'serial_number_hash' => (string) $obit->getSerialNumberHash(),
            'modified_at'        => (string) $obit->getModifiedAt(),
            'obit_status'        => null,
            'owner_did'          => (string) $obit->getOwnerDid(),
            'obd_did'            => (string) $obit->getObdDid(),
            'is_synchronized'    => Obit::NOT_SYNCHRONIZED
        ]);

        $o = Obit::where('obit_did', $obit->getObitId()->toDid())->first();

        $this->assertEquals($o->metadata, $obit->getMetadata()->toArray());
        $this->assertEquals($o->doc_links, $obit->getDocuments()->toArray());
        $this->assertEquals($o->structured_data, $obit->getStructuredData()->toArray());

        $this->assertCount(1, Obit::all());
    }

    /**
     * @test
     */
    public function it_updates_view_record() {
        $this->assertTrue(true);
    }
}
