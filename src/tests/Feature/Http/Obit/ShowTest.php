<?php

declare(strict_types=1);

namespace Tests\Feature\Http\Obit;

use Laravel\Lumen\Testing\DatabaseTransactions;
use Laravel\Lumen\Testing\WithoutEvents;
use Tests\TestCase;
use Illuminate\Support\Facades\Artisan;

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
        Artisan::call('db:seed', ['--class' => 'GenesisSeeder']);

        $genesisDID = 'did:obada:' . sha1('genesis');

        $this->get(route('obits.show', ['obitDID' => $genesisDID]));

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
                'obit_did' => $genesisDID
            ]);
    }
}
