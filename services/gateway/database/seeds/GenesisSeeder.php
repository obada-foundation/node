<?php

use Illuminate\Database\Seeder;
use App\Services\Gateway\ObitDto;
use App\Services\Gateway\Contracts\ServiceContract;
use Carbon\Carbon;

class GenesisSeeder extends Seeder
{
    public function run()
    {
        $obit = new ObitDto([
            'obitDID'    => 'did:obada:' . sha1('genesis'),
            'usn'        => 'none',
            'modifiedAt' => Carbon::now()
        ]);

        app()->make(ServiceContract::class)
            ->create($obit);
    }
}
