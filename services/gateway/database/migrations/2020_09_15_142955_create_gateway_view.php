<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;
use App\Services\Gateway\Models\Obit;

class CreateGatewayView extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('gateway_view', function (Blueprint $table) {
            $table->id();
            $table->unsignedInteger('parent_id')->nullable();
            $table->string('obit_did', 255);
            $table->string('usn', 255);
            $table->json('obit_did_versions')->nullable();
            $table->string('owner_did', 255);
            $table->string('obd_did', 255);
            $table->enum('obit_status', Obit::STATUSES);
            $table->string('manufacturer', 255);
            $table->string('part_number', 255);
            $table->string('serial_number_hash', 255);
            $table->json('metadata')->nullable();
            $table->json('structured_data')->nullable();
            $table->json('doc_links')->nullable();
            $table->dateTime('modified_at');
            $table->string('root_hash', 255);
            $table->tinyInteger('is_synchronized')->default(0);
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('gateway_view');
    }
}
