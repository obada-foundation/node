<?php

declare(strict_types=1);

namespace App\Services\Blockchain;

use App\Services\Blockchain\Contracts\ServiceContract;
use App\Services\Blockchain\QLDB\Driver;

class Service implements ServiceContract {

    protected Driver $driver;

    /**
     * Service constructor.
     * @param Driver $driver
     */
    public function __construct(Driver $driver) {
        $this->driver = $driver;
    }

    /**
     * @param array $obit
     */
    public function create(array $obit) {
        $r = $this->driver->create($obit);

    }

    /**
     * @param string $obitId
     * @return mixed
     */
    public function show(string $obitId)
    {
        // TODO: Implement show() method.
    }

    /**
     * @param string $obitId
     * @return mixed
     */
    public function history(string $obitId)
    {
        // TODO: Implement history() method.
    }

    /**
     * @param string $obitId
     * @return mixed
     */
    public function update(string $obitId)
    {
        $x = $this->ion->decode("224 1 0 234 238 196 129 131 222 192 135 190 189 136 80 101 114 115 111 110 73 100 141 76 105 99 101 110 115 101 78 117 109 98 101 114 139 76 105 99 101 110 115 101 84 121 112 101 141 86 97 108 105 100 70 114 111 109 68 97 116 101 139 86 97 108 105 100 84 111 68 97 116 101 222 186 138 142 150 55 50 85 56 69 56 99 71 88 117 49 74 65 76 99 88 55 49 115 104 71 48 139 139 55 52 52 32 56 52 57 32 51 48 49 140 132 70 117 108 108 141 101 192 15 225 140 134 142 101 192 15 230 138 143");

        $y = $this->ion->decode("224 1 0 234 238 195 129 131 222 192 135 190 189 136 80 101 114 115 111 110 73 100 141 76 105 99 101 110 115 101 78 117 109 98 101 114 139 76 105 99 101 110 115 101 84 121 112 101 141 86 97 108 105 100 70 114 111 109 68 97 116 101 139 86 97 108 105 100 84 111 68 97 116 101 222 186 138 142 150 55 50 85 56 69 56 99 71 88 117 49 74 65 76 99 88 55 49 115 104 71 48 139 139 55 52 52 32 56 52 57 32 51 48 49 140 132 70 117 108 108 141 101 192 15 225 140 134 142 101 192 15 230 138 143");

       /** foreach ($r['ExecuteStatement']['FirstPage']['Values'] as $x) {
            echo "\n\n";
            $arr = unpack("C*", $x['IonBinary']);
            echo implode(" ", $arr);
        }**/
        return response(['x' => $x, 'y' => $y]);
        // TODO: Implement update() method.
    }

    /**
     * @param string $obitId
     * @return mixed
     */
    public function delete(string $obitId)
    {
        // TODO: Implement delete() method.
    }
}
