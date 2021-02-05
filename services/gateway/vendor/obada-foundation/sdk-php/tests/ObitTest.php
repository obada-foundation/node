<?php

declare(strict_types=1);

namespace Tests;

use Obada\Obit;
use Obada\Properties\Document\Document;
use Obada\Properties\Document\HashLink;
use Obada\Properties\Document\Name;
use Obada\Properties\DocumentsCollection;
use Obada\Properties\Manufacturer;
use Obada\Properties\MetadataCollection;
use Obada\Properties\ObitId;
use Obada\Properties\PartNumber;
use Obada\Properties\Property;
use Obada\Properties\SerialNumberHash;
use Obada\Properties\StructuredData\Key;
use Obada\Properties\StructuredData\Record;
use Obada\Properties\StructuredData\Value;
use Obada\Properties\StructuredDataCollection;
use PHPUnit\Framework\TestCase;

class ObitTest extends TestCase {
	public function testItCreatesRootHash(): void {
		$serialNumberHash = hash('sha256', 'SN123456');
		$manufacturer     = 'Sony';
		$partNumber       = 'PN123456';
		$ownerDid         = 'did:obada:owner:123456';
		$modifiedAt       = new \DateTime('now');

		$obitId = new ObitId(
			new SerialNumberHash($serialNumberHash),
			new Manufacturer($manufacturer),
			new PartNumber($partNumber)
		);

		$obit = Obit::make([
			'manufacturer'       => $manufacturer,
			'serial_number_hash' => $serialNumberHash,
			'part_number'        => $partNumber,
			'owner_did'          => $ownerDid,
			'modified_at'        => $modifiedAt,
            'metadata'           => [
                [
                    'key'   => 'type',
                    'value' => 'phone'
                ]
            ],
            'structured_data' => [
                [
                    'key'   => 'color',
                    'value' => 'red'
                ]
            ],
            'documents' => [
                [
                    'name'      => 'swipe report',
                    'hash_link' => 'http://somelink.com'
                ]
            ]
		]);

		$expectedHash = hash('sha256', dechex(
			(new StructuredDataCollection(new Record(new Key('color'), new Value('red'))))->toHash()->toDecimal() +
			(new MetadataCollection(new \Obada\Properties\Metadata\Record(new \Obada\Properties\Metadata\Key('type'), new \Obada\Properties\Metadata\Value('phone'))))->toHash()->toDecimal() +
            (new DocumentsCollection(new Document(new Name('swipe report'), new HashLink('http://somelink.com'))))->toHash()->toDecimal() +
            $obitId->toHash()->toDecimal() +
			hexdec(substr(hash('sha256', $serialNumberHash), 0, 8)) +
			hexdec(substr(hash('sha256', $manufacturer), 0, 8)) +
			hexdec(substr(hash('sha256', $partNumber), 0, 8)) +
			hexdec(substr(hash('sha256', $ownerDid), 0, 8)) +
			hexdec(substr(hash('sha256', (string) $modifiedAt->getTimestamp()), 0, 8)) +
			hexdec(substr(hash('sha256', ''), 0, 8)) + // Obit status
			hexdec(substr(hash('sha256', ''), 0, 8)) // Obd status
		));

		$this->assertEquals($expectedHash, (string) $obit->rootHash());
	}

    /**
     * @dataProvider getGetters
     */
	public function testItCanCallGetMethods(string $method) {
        $obit = Obit::make([
            'manufacturer'       => 'Sony',
            'serial_number_hash' => hash('sha256', 'SN123456'),
            'part_number'        => 'PN123456',
            'owner_did'          => 'did:obada:owner:123456',
            'modified_at'        => new \DateTime('now'),
            'metadata'           => [
                [
                    'key'   => 'type',
                    'value' => 'phone'
                ]
            ],
            'structured_data' => [
                [
                    'key'   => 'color',
                    'value' => 'red'
                ]
            ],
            'documents' => [
                [
                    'name'      => 'swipe report',
                    'hash_link' => 'http://somelink.com'
                ]
            ]
        ]);

        $this->assertInstanceOf(Property::class, $obit->{$method}());
    }

    public function getGetters(): \Generator {
        yield ['getObitId'];
        yield ['getSerialNumberHash'];
        yield ['getManufacturer'];
        yield ['getPartNumber'];
        yield ['getOwnerDid'];
        yield ['getModifiedAt'];
    }

    /**
     * @dataProvider getFailingMethods
     */
    public function testItFailsWhenCallsNotExistingMethod(string $failingMethod) {
	    $this->expectException(\Exception::class);
	    $this->expectExceptionMessage("Method {$failingMethod} is not supported.");

        $obit = Obit::make([
            'manufacturer'       => 'Sony',
            'serial_number_hash' => hash('sha256', 'SN123456'),
            'part_number'        => 'PN123456',
            'owner_did'          => 'did:obada:owner:123456',
            'modified_at'        => new \DateTime('now'),
            'metadata'           => [
                [
                    'key'   => 'type',
                    'value' => 'phone'
                ]
            ],
            'structured_data' => [
                [
                    'key'   => 'color',
                    'value' => 'red'
                ]
            ],
            'documents' => [
                [
                    'name'      => 'swipe report',
                    'hash_link' => 'http://somelink.com'
                ]
            ]
        ]);

        $obit->{$failingMethod}();
    }

    public function getFailingMethods() {
	    yield ['getBar'];
	    yield ['something'];
    }
}