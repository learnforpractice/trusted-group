
import os
import sys
import time
from typing import List, Optional, Dict

try:
	from pyeoskit import eosapi, wallet
except:
	print('pyeoskit not found, please install it with "pip install pyeoskit"')
	sys.exit(-1)

from pyeoskit import log
from pyeoskit.exceptions import ChainException
logger = log.get_logger(__name__)

# modify your test account here
main_account = 'helloworld11'
# modify your test account private key here
wallet.import_key('test', '5JRYimgLBrRLCBAcjHUWCYRv3asNedTYYzVgmiU4q2ZVxMBiJXL')
wallet.import_key('test', '5Jbb4wuwz8MAzTB9FJNmrVYGXo4ABb7wqPVoWGcZ6x8V2FwNeDo')

# {'private': '5JpXLb1tqxJB3Xtzd584xTdqKAzBnQ4TkqfEtT5QPotuv7Yt2bX', 'public': 'PUB_K1_7rVXEPKJsYbjW3KZpaLBnpTx7U3XTdKedEWfPpj2mF9jTixBDr'}
# {'private': '5J2VCfZgiB6g86NfBmJE73yGRpvZXqS2UR7ZDGjt8XU2LPRddiY', 'public': 'PUB_K1_8JEHuDKeLgKfhCQT1rwHSdB1fuDd2ngZDequcvWbF9a8zd9GZC'}
# {'private': '5J3hATRW5G6GT3pkSuyxrPw9KTU4hmbwX9WKddUSLwMcwQSLYex', 'public': 'PUB_K1_6Bi6Ndo1PxiEuvt1osP267xyjo8HEbLgVzBrX7Ewruf7d2MtEk'}
# {'private': '5K2Up4bcBo6BgDpfgccNmpgRE33nrvKsxfMgAauTMzHAqyN8SxM', 'public': 'PUB_K1_8N2voiByjmCZeTwQpjHjvwpg8Gnbv3wFTqgcg7GdKk1UzHqj66'}

wallet.import_key('test', '5JpXLb1tqxJB3Xtzd584xTdqKAzBnQ4TkqfEtT5QPotuv7Yt2bX')
wallet.import_key('test', '5J2VCfZgiB6g86NfBmJE73yGRpvZXqS2UR7ZDGjt8XU2LPRddiY')
wallet.import_key('test', '5J3hATRW5G6GT3pkSuyxrPw9KTU4hmbwX9WKddUSLwMcwQSLYex')
wallet.import_key('test', '5K2Up4bcBo6BgDpfgccNmpgRE33nrvKsxfMgAauTMzHAqyN8SxM')

# EOS7rVXEPKJsYbjW3KZpaLBnpTx7U3XTdKedEWfPpj2mF9jTdcsG5
# EOS8JEHuDKeLgKfhCQT1rwHSdB1fuDd2ngZDequcvWbF9a92VKubc
# EOS6Bi6Ndo1PxiEuvt1osP267xyjo8HEbLgVzBrX7Ewruf7dbUiRp
# EOS8N2voiByjmCZeTwQpjHjvwpg8Gnbv3wFTqgcg7GdKk1UvYekzX

#helloworld12
#'5JHRxntHapUryUetZgWdd3cg6BrpZLMJdqhhXnMaZiiT4qdJPhv',#EOS89jesRgvvnFVuNtLg4rkFXcBg2Qq26wjzppssdHj2a8PSoWMhx
wallet.import_key('test', '5JHRxntHapUryUetZgWdd3cg6BrpZLMJdqhhXnMaZiiT4qdJPhv')

# modify test node here
eosapi.set_node('http://127.0.0.1:9000')

info = eosapi.get_account(main_account)
# logger.info(info)

owner_key = 'EOS7sPDxfw5yx5SZgQcVb57zS1XeSWLNpQKhaGjjy2qe61BrAQ49o'
active_key = 'EOS7sPDxfw5yx5SZgQcVb57zS1XeSWLNpQKhaGjjy2qe61BrAQ49o'
try:
    eosapi.create_account(main_account, 'mtgxinmtgxin', owner_key, active_key, 1024*1024, 1.0, 10000.0)
except Exception as e:
    logger.exception(e)

try:
    eosapi.create_account(main_account, 'mixincrossss', owner_key, active_key, 1024*1024, 1.0, 10000.0)
except Exception as e:
    logger.exception(e)

eosapi.create_account(main_account, 'mixinwtokens', owner_key, active_key, 1024*1024, 1.0, 10000.0)
eosapi.create_account(main_account, 'mixinmanager', owner_key, active_key, 1024*1024, 1.0, 10000.0)

eosapi.create_account(main_account, 'mtgpublisher', owner_key, active_key, 1024*1024, 1.0, 10000.0)
eosapi.create_account(main_account, 'mtgexecutor1', owner_key, active_key, 1024*1024, 1.0, 10000.0)

pub_key = 'EOS7sPDxfw5yx5SZgQcVb57zS1XeSWLNpQKhaGjjy2qe61BrAQ49o'
# pub_key = 'EOS6AjF6hvF7GSuSd4sCgfPKq5uWaXvGM2aQtEUCwmEHygQaqxBSV'

# from pyeoskit import utils
# utils.dbw('helloworld11', 'mtgxinmtgxin', 10.0, 1000.0)

def deploy_contract(account, path):
    with open(f'{path}.wasm', 'rb') as f:
        code = f.read()
    with open(f'{path}.abi', 'rb') as f:
        abi = f.read()
    # if account == 'mtgxinmtgxin':
    #     abi = b''

    try:
        eosapi.deploy_contract(account, code, abi, vm_type=0)
    except ChainException as e:
        if not e.json['error']['details'][0]['message'] == 'contract is already running this version of code':
            raise e

def update_auth(account, parent='owner', permission='active', pub_keys: Optional[List] = None, auth_accounts: Optional[Dict] = None, threshold = 1):
    if not pub_keys and not auth_accounts:
        raise Exception("pub_keys and auth_accounts are empty")

    if not pub_keys:
        pub_keys = []

    if not auth_accounts:
        auth_accounts = {}

    keys = []
    pub_keys.sort()
    for pub_key in pub_keys:
        keys.append({'key': pub_key, 'weight': 1})
    logger.info(keys)

    account_permissions = []
    for a in sorted(auth_accounts.keys()):
        perm = {
            "permission":
            {
                "actor": a,
                "permission": auth_accounts[a],
            },
            "weight":1
        }
        account_permissions.append(perm)

    args = {
        "account": account,
        "permission": permission,
        "parent": parent,
        "auth": {
            "threshold": threshold,
            "keys": keys,
            "accounts": account_permissions,
            "waits": []
        }
    }
    logger.info("++++%s", {account: permission})
    r = eosapi.push_action('eosio', 'updateauth', args, {account: permission})

pub_keys = [
   "EOS7rVXEPKJsYbjW3KZpaLBnpTx7U3XTdKedEWfPpj2mF9jTdcsG5",
   "EOS8JEHuDKeLgKfhCQT1rwHSdB1fuDd2ngZDequcvWbF9a92VKubc",
   "EOS6Bi6Ndo1PxiEuvt1osP267xyjo8HEbLgVzBrX7Ewruf7dbUiRp",
   "EOS8N2voiByjmCZeTwQpjHjvwpg8Gnbv3wFTqgcg7GdKk1UvYekzX"
]

deploy_contract('mtgxinmtgxin', './mtg.xin/mtg.xin')

deploy_contract('mixinwtokens', './mixinproxy/token/token')
deploy_contract('mixincrossss', './mixinproxy/mixinproxy')

update_auth('mixinmanager', pub_keys=pub_keys, threshold=3)
update_auth('mixinmanager', parent='', permission='owner', pub_keys=pub_keys, threshold=3)

# 'mtgxinmtgxin', 'mixincrossss', 'mixinwtokens' 这三个账号的 owner key 都由 mixinmanager 管理
# 'mtgxinmtgxin', 'mixincrossss' 这两个账号的active key都由 mixinmanager 管理
# 'mixinwtokens' 的 active 权限设置成 mixinmanager active

update_auth('mtgxinmtgxin', auth_accounts={'mixinmanager':'active', 'mtgxinmtgxin':'eosio.code'}, threshold=1)
update_auth('mixincrossss', auth_accounts={'mixinmanager':'active', 'mixincrossss':'eosio.code'}, threshold=1)
update_auth('mixinwtokens', auth_accounts={'mixincrossss':'eosio.code', 'mixinwtokens':'eosio.code'}, threshold=1)

update_auth('mtgxinmtgxin', parent='', permission='owner', auth_accounts={'mixinmanager':'active'}, threshold=1)
update_auth('mixincrossss', parent='', permission='owner', auth_accounts={'mixinmanager':'active'}, threshold=1)
update_auth('mixinwtokens', parent='', permission='owner', auth_accounts={'mixinmanager':'active'}, threshold=1)
