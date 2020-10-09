<?php

declare(strict_types=1);

namespace Tests;

use Laravel\Lumen\Testing\TestCase as BaseTestCase;

use App\Services\Gateway\ObitDto;
use Faker\Generator as Faker;
use Carbon\Carbon;
use App\Obada\ObitId;

abstract class TestCase extends BaseTestCase
{
    /**
     * Creates the application.
     *
     * @return \Laravel\Lumen\Application
     */
    public function createApplication()
    {
        return require __DIR__.'/../bootstrap/app.php';
    }

    public function validObitDto() {
        $faker = new Faker;
        $faker->addProvider(new \Faker\Provider\DateTime($faker));
        $faker->addProvider(new \Faker\Provider\Lorem($faker));
        $faker->addProvider(new \Faker\Provider\en_US\Person($faker));
        $faker->addProvider(new \Faker\Provider\Company($faker));

        $manufacturer = $faker->company;
        $hash         = hash('sha256', $faker->word);
        $partNumber   = $faker->word;
        $obit         = new ObitId($hash, $manufacturer, $partNumber);

        try {
            return new ObitDto([
                'serialNumberHash' => $hash,
                'modifiedAt'       => Carbon::parse($faker->dateTime->getTimestamp())->format('Y-m-d H:i:s'),
                'manufacturer'     => $manufacturer,
                'usn'              => $obit->toUsn(),
                'partNumber'       => $partNumber,
                'obitDID'          => $obit->toHash()
            ]);
        } catch (\Throwable $t) {
            dd($t->getMessage());
        }
    }
}
