<?php

declare(strict_types=1);

namespace Tests\Feature\Http\Obit;

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
        $did = 'did:obada:owner:123456';

        $this->get(route('obits.show', ['obitDID' => $did]));

        $this->seeStatusCode(200)
            ->seeJsonStructure([
                'code',
                'message'
            ])
            ->seeJson([
                'code'    => 404,
                'message' => "Obit with did \"{$did}\" doesn't exists"
            ]);
    }
}
