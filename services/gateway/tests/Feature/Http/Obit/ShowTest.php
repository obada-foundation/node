<?php

declare(strict_types=1);

namespace Tests\Feature\Http\Obit;

use App\Services\Gateway\Models\Obit;
use Laravel\Lumen\Testing\DatabaseTransactions;
use Laravel\Lumen\Testing\WithoutEvents;
use Tests\TestCase;

class ShowTest extends TestCase {

    use DatabaseTransactions, WithoutEvents;

    /**
     * @test
     */
    public function it_returns_404_when_obit_not_exists() {
        $did = 'did:obada:owner:123456';

        $this->get(route('obits.show', ['obitDID' => $did]));

        $this->seeStatusCode(404)
            ->seeJsonStructure([
                'code',
                'message'
            ])
            ->seeJson([
                'code'    => 404,
                'message' => "Obit with did \"{$did}\" doesn't exists"
            ]);
    }

    /**
     * @test
     */
    public function it_returns_a_genesis_record() {
        $obit = factory(Obit::class)->create();

        $this->get(route('obits.show', ['obitDID' => $obit->obit_did]));

        $this->seeStatusCode(200)
            ->seeJsonStructure([
                'doc_links' => [],
                'is_synchronized',
                'manufacturer',
                'metadata' => [],
                'modified_at',
                'obd_did',
                'obit_did',
                'obit_did_versions' => [],
                'obit_status',
                'owner_did',
                'part_number',
                'qldb_root_hash',
                'root_hash',
                'serial_number_hash',
                'structured_data' => [],
                'usn'
            ])
            ->seeJson([
                'obit_did' => $obit->obit_did
            ]);
    }
}
